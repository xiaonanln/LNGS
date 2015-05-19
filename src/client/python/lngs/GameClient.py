# coding: utf8

import socket
from RPCMessenger import RPCMessenger
import EntityManager
import logging 

class GameClient(object):
	def __init__(self, host, port):
		self.sock = socket.socket()
		self.sock.connect((host, port))

		self.rpc = RPCMessenger(self.sock)
		self.rpc.OnRecvMessage = self._on_recv_message

		self._entity_classes = {}

	def _on_recv_message(self, msg):
		print 'recv', msg 
		if 'M' in msg:
			self._handle_CallMethod(msg['ID'], msg['M'], msg.get('ARGS', ()))
		if 'CE' in msg:
			# CreateEntity message
			self._handle_CreateEntity(*msg['CE'])
		elif 'DE' in msg:
			self._handle_DestroyEntity(msg['DE'])

	def _handle_CallMethod(self, entityid, methodname, args):
		print 'CallMethod', entityid, methodname, args
		entity = EntityManager.getentity(entityid)
		if not entity:
			logging.error('CallMethod %s error: entity not found: %s', methodname, entityid)
			return 

		method = getattr(entity, methodname)
		logging.info('calling entity method %s.%s', entity, method)
		method(*args)

	def _handle_CreateEntity(self, entityid, entitytype):
		print 'CreateEntity', entityid, entitytype
		print self._entity_classes
		clz = self._entity_classes[entitytype]
		print clz
		entity = clz(entityid)
		print clz, entity
		EntityManager.addentity(entity)

	def _handle_DestroyEntity(self, entityid):
		print 'DestroyEntity', entityid
		EntityManager.delentity(entityid)

	def call_entity(self, entityid, method, args):
		msg = {
			'ID': entityid, 
			'M' : method, 
			'ARGS': args
		}

		self.rpc.send_message(msg)

	def register_entity_class(self, entity_clz):
		print 'register_entity_class', entity_clz.__name__
		self._entity_classes[entity_clz.__name__] = entity_clz

