#!/usr/bin/env python

import json
import sys

for line in sys.stdin:
    decoded = json.loads(line)
    print(json.dumps(decoded, indent=4))
