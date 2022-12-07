#!/usr/bin/env python

import boilerplate
import unittest
from io import StringIO
import os
import sys

class TestBoilerplate(unittest.TestCase):
  """
  Note: run this test from the hack/boilerplate directory.

  $ python -m unittest boilerplate_test
  """

  def test_boilerplate(self):
    os.chdir("test/")

    class Args(object):
      def __init__(self):
        self.filenames = []
        self.rootdir = "."
        self.boilerplate_dir = "../"
        self.verbose = True

    # capture stdout
    old_stdout = sys.stdout
    sys.stdout = StringIO.StringIO()

    boilerplate.args = Args()
    ret = boilerplate.main()

    output = sorted(sys.stdout.getvalue().split())

    sys.stdout = old_stdout

    self.assertEquals(
        output, ['././fail.go', '././fail.py'])
