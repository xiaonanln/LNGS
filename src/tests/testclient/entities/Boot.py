#  coding: utf8

import lngs 
from lngs import Timer
import logging 

class Boot(lngs.Entity): 
	def BecomePlayer(self):
		logging.info('BecomePlayer %s', self)

		Timer.addRepeatTimer(1, self.playGame)
		
	def playGame(self):
		lngs.client.call_entity(self.id, 'PlayGame', [100, "test string", {'test': 'dict'}, ['test', 'list']])