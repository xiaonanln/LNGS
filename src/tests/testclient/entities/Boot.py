#  coding: utf8

import lngs 
from lngs import Timer
import logging 

USERNAME, PASSWORD = "123456", "123456"

class Boot(lngs.Entity): 
	def BecomePlayer(self):
		logging.info('BecomePlayer %s', self)

		Timer.addTimer(0.1, lambda: self.Register() )
		
	def playGame(self):
		lngs.client.call_entity(self.id, 'PlayGame', [100, "test string", {'test': 'dict'}, ['test', 'list']])

	def Login(self, username, password):
		lngs.client.call_entity(self.id, 'Login', [username, password])

	def OnLogin(self, result, username):
		logging.info('OnLogin %s %s', result, username)

	def Register(self):
		lngs.client.call_entity(self.id, "Register", [USERNAME, PASSWORD])
		
	def OnRegister(self, result):
		logging.info("OnRegister %s", result)
		self.Login( USERNAME, PASSWORD ) 