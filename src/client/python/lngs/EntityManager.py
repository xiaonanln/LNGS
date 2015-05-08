# coding: utf8

class EntityManager(object):

	def __init__(self):
		self.entities = {}

	def addentity(self, entity):
		self.entities[entity.id] = entity 

	def delentity(self, entity):
		try: 
			del self.entities[entity.id]
		except:
			pass 

