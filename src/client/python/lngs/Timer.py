# coding: utf8

import time 
import heapq

_timerQueue = []

def addTimer(timeout, callback):
	timer = Timer(timeout, callback, False)
	heapq.heappush(_timerQueue, (timer.firetime, timer))

def addRepeatTimer(timeout, callback):
	timer = Timer(timeout, callback, True)
	heapq.heappush(_timerQueue, (timer.firetime, timer))

def _loop():
	now = time.time()
	while len(_timerQueue) > 0:
		firetime, timer = _timerQueue[0]
		if now >= firetime:
			# fire the callback
			timer.fire()
			heapq.heappop(_timerQueue)
		else:
			break 

class Timer(object):
	def __init__(self, timeout, callback, repeat=False):
		self.timeout = timeout
		self.firetime = time.time() + timeout
		self.callback = callback
		self.repeat = repeat 

	def fire(self):
		try:
			self.callback()
		except:
			import traceback
			traceback.print_exc()

		if self.repeat:
			self.firetime = time.time() + self.timeout # next fire time
			heapq.heappush(_timerQueue, (self.firetime, self))




