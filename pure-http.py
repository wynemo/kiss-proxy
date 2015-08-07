#!/usr/bin/env python
# coding:utf-8

"""
a simple and stupid http proxy server
"""

import socket
import select
import logging
import urlparse
from SocketServer import ThreadingTCPServer, StreamRequestHandler


CRLF = '\r\n'


class Encoder(StreamRequestHandler):
    def handle_tcp(self, sock, remote, first_data):
        fdset = [sock, remote]

        if first_data:
            if remote.send(first_data) <= 0:
                return

        while True:
            r, w, e = select.select(fdset, [], [])
            try:
                if sock in r:
                    data = sock.recv(64 * 1024)
                    logging.debug('data from browser ' + str(len(data)))
                    if remote.send(data) <= 0:
                        break
                if remote in r:
                    data = remote.recv(64 * 1024)
                    logging.debug('data from proxy ' + str(len(data)))
                    if sock.send(data) <= 0:
                        break
            except socket.sslerror as e:
                if e.args[0] == socket.SSL_ERROR_EOF:
                    break
                else:
                    raise

    def handle(self):
        remote = None
        try:
            sock = self.connection
            first_line = ''
            pos = 0
            data = sock.recv(64 * 1024)
            if not data:
                raise Exception('hey no data, don\'t  fuck with me')
            logging.debug('socks connection from ' + str(self.client_address))
            while not first_line:
                pos = data.find(CRLF, 0)
                if pos != -1:
                    first_line = data[:pos]
                else:
                    data += sock.recv(64 * 1024)
            old_pos = pos + 2
            first_line_end = old_pos

            method, url, version = first_line.strip().split(' ')
            import re

            if method.upper() != 'CONNECT':
                ss = urlparse.urlsplit(url)
                path = ss.path
                port = ss.port or 80
                if ss.query:
                    path += '?' + ss.query
                if ss.fragment:
                    path += '#' + ss.fragment
                if not path:
                    path = '/'
                host = ss.hostname
            else:
                host, port = url.split(':')
                port = int(port)
                path = ''

            modified_first_line = ' '.join([method, path, version])
            modified_data = modified_first_line + CRLF + data[first_line_end:]
            modified_data = modified_data.replace('Proxy-Connection: keep-alive\r\n', 'Connection: keep-alive\r\n', 1)
            remote = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            if method.upper() == 'CONNECT':
                resp = CRLF.join(['HTTP/1.1 200 Connection established',
                                  'Proxy-agent: proxy.py v1',
                                  CRLF])
                sock.send(resp)
                modified_data = ''
            l = (host, port)
            remote.connect(l)
            self.handle_tcp(sock, remote, modified_data)
        except Exception as e:
            logging.exception(first_line + ' ' + str(e))
        finally:
            if remote:
                remote.close()


def main():
    level = logging.INFO
    logging.basicConfig(format='%(asctime)s [%(levelname)s] %(message)s',
                        datefmt='%m/%d/%Y %I:%M:%S %p',
                        level=level)
    ThreadingTCPServer.allow_reuse_address = True
    server = ThreadingTCPServer(('0.0.0.0', 8118), Encoder)
    server.serve_forever()


if __name__ == '__main__':
    main()
