#!/usr/bin/env python

import argparse
import glob
import os
import sys
import random
import atexit

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
parser = argparse.ArgumentParser(
    prog="termtosvg.py", description="wrapper to get some terminal screenshot"
)
parser.add_argument("output")
parser.add_argument("--command", type=str, default="date")
parser.add_argument("--lines", type=int, default=10)
parser.add_argument("--columns", type=int, default=100)

args = parser.parse_args()

rndint = random.randint(0, 1000000)
tmpdir = f"termtosvg_{rndint}"
atexit.register(lambda: os.system(f"rm -Rf {tmpdir}"))

cmd = f"termtosvg --still-frames --screen-geometry '{args.columns}x{args.lines}' --template=solarized_light --command '{args.command}' '{tmpdir}'"
print(cmd)
rc = os.system(cmd)
if rc != 0:
    print("ERROR during termtosvh, output:")
    sys.exit(1)

svg = sorted(glob.glob(os.path.join(tmpdir, "*.svg")))[-1]
os.system("cp -f %s %s" % (svg, args.output))

print()
print(f"=> {args.output} is ready!")
