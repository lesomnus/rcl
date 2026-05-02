#!/usr/bin/env python3
import sys
import rclpy
from rclpy.node import Node
from example_interfaces.srv import AddTwoInts

topic = sys.argv[1] if len(sys.argv) > 1 else 'add_two_ints'

class S(Node):
    def __init__(self):
        super().__init__('s')
        self.create_service(AddTwoInts, topic, self.cb)

    def cb(self, req, res):
        res.sum = req.a + req.b
        return res

rclpy.init()
rclpy.spin(S())
