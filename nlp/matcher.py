import json
import spacy
from spacy.matcher import Matcher
from spacy.matcher import PhraseMatcher

def is_free_shipping(doc):
    matcher = Matcher(nlp.vocab)
    matcher.add("target", None,[{"LOWER": "free shipping"}])
    matches = matcher(doc)
    return len(matches) > 0

def extract_free_item(matcher, doc):
    # in the format of "NUM free product"
    matcher.add("target", None, [{"IS_DIGIT": True, "OP": "?"}, {"LOWER": "free"}, {"ENT_TYPE": "CATALOG"}])
    matches = matcher(doc)
    for match_id, start, end in matches:
        #  print("found free item:", doc[start:end])
        if doc[start].is_digit:
            return docs[start].text, doc[start+2,end]
        else:
            # quantity is omitted to 1 when not specified explicitly
            return "1", doc[start+1:end]

    # in the format of "get NUM product free"
    #  matcher.add("target", None, [{"IS_DIGIT": True, "OP": "?"}, {"IS_ASCII": True, "OP": "+"}, {"LOWER": "free"}])
    #  matches = matcher(doc)
    #  for match_id, start, end in matches:
        #  print("found free item", doc[start:end].text)
        #  return 1, doc[start:end-1].text

def extract_discount_target(matcher, doc):
    matcher.add("target", None, [{"LOWER": "off"}, {"IS_ALPHA": True, "OP": "*"}, {"DEP": "pobj"}])
    matches = matcher(doc)
    for match_id, start, end in matches:
        return doc[end-1:end]

def find_money_target(span):
    for token in span:
        if token.ent_type_ == "CATALOG":
            return token
        if token.ent_type_ in ["MONEY", "PERCENTAGE"]:
            # we are running into another condition/action part, stop now
            break
    return None # use None to signal the target is order


def extract_money(doc):
    all_money = []
    condition_verbs = ["spend", "purchase", "buy"]
    condition_adps = ["over", "with"]
    action_verbs = ["get", "reward", "receive"]
    preceding_verb = None
    preceding_adp = None
    for token in doc:
        if token.ent_type_ == "MONEY" and token.is_digit:
            # look to the right to see if this is a $X off case
            if doc[token.i+1].lower_ == "off":
                target = find_money_target(doc[token.i+1:])
                all_money.append((token.text, "action", target))
            # look to the left to see if this is a spend|get $X case
            elif preceding_verb in condition_verbs or preceding_adp in condition_adps:
                target = find_money_target(doc[token.i+1:])
                all_money.append((token.text, "condition", target))
            elif preceding_verb in action_verbs:
                target = find_money_target(doc[token.i+1:])
                all_money.append((token.text, "action", target))
        elif token.pos_ == "VERB":
            preceding_verb = token.lower_
        elif token.pos_ == "ADP":
            preceding_adp = token.lower_

    return all_money

def is_action(money):
    value, pos, target = money
    return pos == "action"

def is_condition(money):
    value, pos, target = money
    return pos == "condition"

def extract_percentage(doc):
    for token in doc:
        if token.ent_type_ == "PERCENT" and token.is_digit:
            target = find_money_target(doc[token.i+2:])
            return (token.text, target)

def map_catalog_token_to_api(token):
    if token.text in products:
        key = "products"
    elif token.text in brands:
        key = "brands"
    else:
        key = "categories"
    return {key: [token.text]}

def create_free_shipping(doc):
    if is_free_shipping(doc):
        return {"shipping": {"free_shipping": True, "zone_ids": "*"}}

    return None

def create_free_item(doc, catalog_indexes, catalog_qty):
    action_index = index_free_item_action_token(doc)
    if action_index is None:
        return None

    free_item_index = closest_index(action_index, catalog_indexes)
    if free_item_index is None:
        return None

    free_item = doc[free_item_index]
    # qty has matching list index per catalog value
    free_qty = catalog_qty[catalog_indexes.index(free_item_index)]

    return {"cart_items": {"discount": {"percentage_amount": 100}, "items": map_catalog_token_to_api(free_item), "quantity": free_qty, "add_free_item": True}}

def create_percentage_discount(doc, order_indexes, catalog_indexes, catalog_qty):
    for pi in index_percent_token(doc):
        closest_order = closest_index(pi, order_indexes)
        closest_catalog = closest_index(pi, catalog_indexes)
        if closest_order < closest_catalog:
            # percentage discount on order
            return {"cart_value": {"discount": {"percentage_amount": int(doc[pi])}}}
        else:
            # percentage discount on catalog
            items = map_catalog_token_to_api(doc[closest_catalog])
            qty = catalog_qty[catalog_indexes.index(closest_catalog)]
            return {"cart_items": {"discount": {"percentage_amount": int(doc[pi])}, "items": items, "quantity": qty}}

def create_fixed_discount(doc, money_indexes, order_indexes, catalog_indexes, catalog_qty):
    for i in find_action_indexes(doc):
        money_index = closest_index(i, money_indexes)
        closest_order = closest_index(money_index, order_indexes)
        closest_catalog = closest_index(money_index, catalog_indexes)
        if closest_order < closest_catalog:
            # percentage discount on order
            return {"cart_value": {"discount": {"fixed_amount": int(doc[money_index])}}}
        else:
            # percentage discount on catalog
            items = map_catalog_token_to_api(doc[closest_catalog])
            qty = catalog_qty[catalog_indexes.index(closest_catalog)]
            return {"cart_items": {"discount": {"fixed_amount": int(doc[money_index])}, "items": items, "quantity": qty}}

def createAction(doc):
    result = extract_free_item(Matcher(nlp.vocab), doc)
    if result is not None:
        qty, free_item_span = result
        for item_token in free_item_span:
            item = map_catalog_token_to_api(item_token)

        return {"cart_items": {"discount": {"percentage_amount": 100}, "items": item, "quantity": qty, "add_free_item": True}}

    fixed_amount = None
    percentage_amount = None
    target = None

    money_values = extract_money(doc)
    discount_moneys = filter(is_action, money_values)
    for money in discount_moneys:
        fixed_amount, pos, target = money

    percentage = extract_percentage(doc)
    if percentage is not None:
        percentage_amount, target = percentage

    if fixed_amount is not None:
        discount = {"fixed_amount": fixed_amount}
    elif percentage_amount is not None:
        discount = {"percentage_amount": percentage_amount}
    else:
        print("no discount found for this line, skip")
        return None

    if target is None:
        action = {"cart_value": {"discount": discount}}
    else:
        item = map_catalog_token_to_api(target)
        action = {"cart_items": {"discount": discount, "items": item, "quantity": "1"}}

    return action

#  def createCondition(doc):
    #  money_values = extract_money(doc)
    #  condition_moneys = filter(is_condition, money_values)
    #  fixed_amount = None
    #  for money in condition_moneys:
        #  fixed_amount, pos, target = money
        #  return{"cart": {"minimum_spend": fixed_amount}}

def merge_token(matcher, doc, id, matches, entity_type):
    with doc.retokenize() as retokenizer:
        for match_id, start, end in matches:
            retokenizer.merge(doc[start:end], attrs={"ENT_TYPE": entity_type})

def merge_catalog_token(matcher, doc, id, matches):
    merge_token(matcher, doc, id, matches, "CATALOG")

# match and merge all the catalog phrase into single token for easy token matching later
def mark_catalog(doc, terms):
    matcher = PhraseMatcher(nlp.vocab, attr="LOWER")
    patterns = [nlp.make_doc(text) for text in terms]
    matcher.add("catalog", merge_catalog_token, *patterns)
    matches = matcher(doc)

def index_catalog_token(doc):
    pos = []
    for token in doc:
        if token.ent_type_ == "CATALOG":
            pos.append(token.i)

    return pos

def find_catalog_qty(doc, catalog_indexes):
    matching_qty = [];
    for index in catalog_indexes:
        if index > 1 and doc[index-1].ent_type_ == "CARDINAL" and doc[index-1].pos_ == "NUM":
            matching_qty.append(int(doc[index-1]))
        else:
            matching_qty.append(1) # default qty to 1 when not found

    return matching_qty

def index_money_token(doc):
    pos = []
    for token in doc:
        if token.ent_type_ == "MONEY" and token.is_digit:
            pos.append(token.i)

    return pos

def index_percent_token(doc):
    pos = []
    for token in doc:
        if token.ent_type_ == "PERCENT" and token.is_digit:
            pos.append(token.i)

    return pos

def mark_free_item_action(doc):
    matcher = PhraseMatcher(nlp.vocab, attr="LOWER")
    patterns = [nlp.make_doc(text) for text in ["free"]
    matcher.add("free_item", merge_free_item_token, *patterns)
    matches = matcher(doc)

def merge_free_item_token(matcher, doc, id, matches):
    merge_token(matcher, doc, id, matches, "ACTION")

def find_action_indexes(doc):
    matcher = Matcher(nlp.vocab)
    matcher.add("discount_item", None, [{"LOWER": "get"}], [{"LOWER": "off"}]])
    matches = matcher(doc)
    indexes = []
    for match_id, start, end in matches:
        indexes.append(start)

    return indexes

def merge_discount_item_token(matcher, doc, id, matches):
    merge_token(matcher, doc, id, matches, "ACTION")

def index_free_item_action_token(doc):
    for token in doc:
        if token.lower_ == "free":
            return token.i
    return None

def mark_item_condition(doc):
    matcher = PhraseMatcher(nlp.vocab, attr="LOWER")
    terms = ["buy", "purchage"]
    patterns = [nlp.make_doc(text) for text in terms]
    matcher.add("free_action", merge_item_condition_token, *patterns)
    matches = matcher(doc)

def merge_item_condition_token(matcher, doc, id, matches):
    merge_token(matcher, doc, id, matches, "CONDITION")

def index_item_condition_token(doc):
    pos = []
    for token in doc:
        if token.ent_type_ == "CONDITION":
            pos.append(token.i)
    return pos

def index_order_token(doc):
    matcher = Matcher(nlp.vocab)
    matcher.add("target", None,[{"LOWER": "order"}])
    matches = matcher(doc)
    indexes = []
    for match_id, start, end in matches:
        indexes.append(start)

    return indexes

def closest_index(anchor, indexes):
    if len(indexes) == 0:
        return None

    shortest = 100 # randomly picked large enough number for comparison purpose
    closest = anchor
    for i in indexes:
        distance = abs(anchor - i)
        if distance <= shortest:
            shortest = distance # favor the one on the right with eaqual distance
            closest = i

    return closest

#  print([(token.text, token.ent_type_) for token in doc])
nlp = spacy.load('en_core_web_sm')

text = """
Buy a pair of earrings and get a second pair for free (Randem)
BUY 2 TEES, GET THE 3RD FREE (superdry)
"""

# TODO: need to sync from store catalog to populate these data
products = ["Rentinol 24", "3-piece set", "Alphard Club Booster eWheels"]
categories = ["drivers", "tees", "earrings"]
brands = ["Puma"]

for line in text.splitlines():
    doc = nlp(line)

    # mark catalog items by merging token with ent_type=CATALOG
    catalog = products.copy()
    catalog.extend(categories)
    catalog.extend(brands)
    mark_catalog(doc, catalog)
    catalog_indexes = index_catalog_token(doc)
    catalog_qty = find_catalog_qty(doc, catalog_indexes)

    order_indexes = index_order_token(doc)
    money_indexes = index_money_token(doc)

    action = create_free_shipping(doc) or \
             create_free_item(doc, catalog_indexes, catalog_qty) or \
             create_percentage_discount(doc, order_indexes, catalog_indexes, catalog_qty) or \
             create_fixed_discount(doc, money_indexes, order_indexes, catalog_indexes, catalog_qty)

    #  mark_item_condition(doc)

    item_condition_indexes = index_item_condition_token(doc)

    # given if item is specified in condition, it's sometimes omitted from action
    # here we try to identify an action and then search for the closest item for it
    free_item_index = None
    free_item = None
    action = None

    item_condition_index = None
    item_condition = None
    condition = None
    for i in item_condition_indexes:
        item_condition_index = closest_index(i, catalog_indexes)
        break

    if item_condition_index is not None:
        condition_item = doc[item_condition_index]
        # for now omit free quantity to be 1
        condition = {"cart": {"items": map_catalog_token_to_api(condition_item), "minimum_quantity": 1}}

    #  print([(token.text, token.ent_type_) for token in doc])
    #  action = createAction(doc)
    #  if action is None:
        #  continue

    if action is None:
        print("no action for this rule, skip")
        continue

    rule = {"action": action}

    if condition is not None:
        rule["condition"] = condition

    #  condition = createCondition(doc)
    #  if condition is not None:
        #  rule["condition"] = condition

    print(line, "=>", json.dumps(rule))
