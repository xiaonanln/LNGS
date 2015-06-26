#  coding: utf8

import lngs 
from lngs import Timer
import logging 

class Avatar(lngs.Entity): 
	def BecomePlayer(self):
		logging.info('BecomePlayer %s', self)
		Timer.addRepeatTimer(1, lambda: self.AddExp(10))
	
	def AddExp(self, exp):
		logging.info('AddExp %d', exp)
		lngs.client.call_entity(self.id, 'AddExp', [10])
	