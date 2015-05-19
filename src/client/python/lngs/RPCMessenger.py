# coding: utf8


import sys 
from asyncore import dispatcher
import traceback
from bson import BSON
from bson.objectid import ObjectId

import struct
import Timer

class InvalidPacketError(Exception): pass

class Packeter(object):
	
	MAX_PACKET_LENGTH = 1 * 1024 * 1024
	
	def  __init__(self):
		self.buffer = ''
		self.packet_length = -1
	
	def feed(self, data):
		self.buffer += data
		packets = []
		
		while True:
			if self.packet_length == -1:
				if len(self.buffer) >= 4:
					# read packet length
					self.read_packet_length()
				else:
					break 
			else:
				if len(self.buffer) >= self.packet_length:
					# get the packet data
					pd, self.buffer = self.buffer[:self.packet_length], self.buffer[self.packet_length:]
					packets.append(pd)
					self.packet_length = -1
				else:
					break
				
		return packets
		
	def read_packet_length(self):
		lendata, self.buffer = self.buffer[:4], self.buffer[4:]
		self.packet_length = struct.unpack('<I', lendata)[0]
		if self.packet_length > Packeter.MAX_PACKET_LENGTH:
			raise InvalidPacketError("invalid packet length: %d" % self.packet_length)
		

class RPCMessenger(dispatcher):
	
	def __init__(self, sock):
		dispatcher.__init__(self, sock)
		self.packeter = Packeter()
		self.out_buffer = ''
		self.OnRecvMessage = None
		
	def handle_read(self):
		data = self.recv(8192)
		packets = self.packeter.feed(data)
		print 'handle_read', len(packets), 'packets'
		if self.OnRecvMessage:
			for packet in packets:
				try:
					msg = BSON(packet).decode()
					self.OnRecvMessage(msg)
				except Exception, e:
					print >>sys.stderr, 'OnRecvMessage error: %s' % e

		print 'handle_read done'

	def handle_close(self):
		print >>sys.stderr, "TcpClient closed: %s, %s" % (self,self.connected)
		self.close()
		
	def __str__(self):
		return "%s:%d" % self.addr
		
	def initiate_send(self):
		num_sent = dispatcher.send(self, self.out_buffer[:8192])
		self.out_buffer = self.out_buffer[num_sent:]

	def handle_write(self):
		self.initiate_send()

	def readable(self):
		return True

	def writable(self):
		Timer._loop()
		return len(self.out_buffer) > 0

	def send(self, data):
		self.out_buffer = self.out_buffer + data
		self.initiate_send()

	def send_message(self, msg):
		try:
			payload = BSON.encode(msg)
			length = struct.pack("=I", len(payload))
			self.send( length + payload )
		except Exception, e:
			print >>sys.stderr, 'send_message error: %s: %s' % (e, msg)

