# coding: utf8

from Entity import Entity

entities = {}

def _getid(entity_or_id):
	return entity_or_id.id if isinstance(entity_or_id, Entity) else entity_or_id

def addentity(entity):
	print 'EntityManager.addentity', entity
	entities[entity.id] = entity 

def delentity(entity):
	print 'EntityManager.delentity', entity
	try: del entities[_getid(entity)] 
	except: pass 

def getentity(entityid):
	return entities.get( _getid(entityid) )
