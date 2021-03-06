# coding: utf8

import sys
import os
import logging 

logging.getLogger().setLevel(logging.DEBUG)

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

lngs.connect('127.0.0.1', 7777)

import entities
from entities.Boot import Boot
from entities.Avatar import Avatar
lngs.client.register_entity_class(Boot)
lngs.client.register_entity_class(Avatar)

try:
	username, password = sys.argv[1:3]
except:
	username, password = "123456", "123456"
	
logging.debug("test client running with username = %s, password = %s", username, password)
entities.Boot.USERNAME = username
entities.Boot.PASSWORD = password

while True:
	lngs.loop()

print "This is test client of LNGS"


