#  coding: utf8

import lngs 
from lngs import Timer
import logging 

class Avatar(lngs.Entity): 
	def BecomePlayer(self):
		logging.info('BecomePlayer %s', self)
		lngs.client.call_entity(self.id, 'AddExp', [10])