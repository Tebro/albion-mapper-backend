import sys
import json


def handle_json(data):
  cluster = data["world"]["clusters"]["cluster"]
  selected_fields = [{"name": m["@displayname"], "type": m["@type"]} for m in cluster]

  filtered_by_type = [o for o in selected_fields if (o["type"].startswith("OPENPVP") or o["type"].startswith("TUNNEL") or o["type"] == "SAFEAREA")]

  print(json.dumps(filtered_by_type))

  #distinct_types = set([o["type"] for o in filtered_by_type])
  #print(distinct_types)

  #name_to_find = "Xerites-Oxoulum"
  #name_to_find = "Whitebank Cross"
  #name_to_find = "Qiient-Qinsum"
  #name_to_find = "Wyre Forest"
  #name_to_find = "Oakcopse"
  #print([o for o in selected_fields if o["name"] == name_to_find])



file_path = sys.argv[1]

with open(file_path) as file:
  raw = file.read()
  parsed = json.loads(raw)

  handle_json(parsed)