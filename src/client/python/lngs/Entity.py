# coding: utf8

from bson.objectid import ObjectId

class Entity(object):
	def __init__(self, id=None):
		if id is None:
			id = self._new_id()

		self.id = id 

	@staticmethod
	def _new_id():
		return str(ObjectId())

	def __str__(self):
		return 'entity<%s>' % self.id
