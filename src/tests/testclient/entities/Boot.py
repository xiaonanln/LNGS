#  coding: utf8

import lngs 
from lngs import Timer
import logging 

class Boot(lngs.Entity): 
	def BecomePlayer(self):
		logging.info('BecomePlayer %s', self)

		Timer.addTimer(0.1, self.Register)
		
	def playGame(self):
		lngs.client.call_entity(self.id, 'PlayGame', [100, "test string", {'test': 'dict'}, ['test', 'list']])

	def login(self):
		lngs.client.call_entity(self.id, 'Login', ['test', '123456'])

	def Register(self):
		lngs.client.call_entity(self.id, "Register", ["test3", "123456"])
		
	def OnRegister(self, result):
		logging.info("OnRegister %s", result)