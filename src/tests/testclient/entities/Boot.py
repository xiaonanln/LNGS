#  coding: utf8

import lngs 
import logging 

class Boot(lngs.Entity): 
	def BecomePlayer(self):
		logging.info('BecomePlayer %s', self)
