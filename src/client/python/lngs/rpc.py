#coding: utf8

import asyncore, socket
import logging

class RPCMessenger(asyncore.dispatcher):

	DISCONNECTED = 0
	CONNECTING = 1
	CONNECTED = 2

	def __init__(self, host, port):
		asyncore.dispatcher.__init__(self)
		self.create_socket(socket.AF_INET, socket.SOCK_STREAM)
		self.state = RPCMessenger.CONNECTING
		self.connect( (host, port) )

	def handle_connect(self):
		logging.debug('handle_connect')

	def handle_close(self):
		self.close()

	def handle_read(self):
		print self.recv(8192)

	def writable(self):
		return (len(self.buffer) > 0)

	def handle_write(self):
		sent = self.send(self.buffer)
		self.buffer = self.buffer[sent:]