# coding: utf8

import sys
import os

def setup_client_lib_path():
	exe = sys.argv[0]
	exe_dir = os.path.dirname(exe)
	if exe_dir == '':
		exe_dir = '.'

	src_path = os.path.abspath(os.path.join(exe_dir, '..', '..'))
	client_lib_path = os.path.join(src_path, 'client', 'python')
	print 'Add client lib path: %s' % client_lib_path
	sys.path.append(client_lib_path)

setup_client_lib_path()

import lngs

while True:
	lngs.loop()

print "This is test client of LNGS"


