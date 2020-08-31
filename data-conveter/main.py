import sys
import json

def extract_map_markers(m):
  if "minimapmarkers" in m.keys():
    markers = m["minimapmarkers"]["marker"]
    if markers is None:
      return []
    elif type(markers) is dict:
      return [markers]
    else:
      return markers
  else:
    return []

def extract_resources(m):
  if "distribution" in m.keys() and "resource" in m["distribution"].keys():
    resources = m["distribution"]["resource"]
    if resources is None:
      return []
    elif type(resources) is dict:
      return [resources]
    else:
      return resources
  else:
    return []

def handle_json(data):
  cluster = data["world"]["clusters"]["cluster"]
  selected_fields = [{
      "name": m["@displayname"],
      "type": m["@type"],
      "resources": extract_resources(m),
      "markers": extract_map_markers(m)
  } for m in cluster]

  filtered_by_type = [o for o in selected_fields if (o["type"].startswith("OPENPVP") or o["type"].startswith("TUNNEL") or o["type"] == "SAFEAREA")]

  print(json.dumps(filtered_by_type))



file_path = sys.argv[1]

with open(file_path) as file:
  raw = file.read()
  parsed = json.loads(raw)

  handle_json(parsed)
