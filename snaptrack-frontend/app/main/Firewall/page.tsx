"use client";

import React, { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { FaStop, FaPlus, FaTimes } from "react-icons/fa";
import { useSocket } from "@/app/context/SocketContext";
import { toast } from "react-toastify";

interface FirewallRule {
  id: string;
  protocol: string;
  port: string;
  source: string;
  destination: string;
  action: string;
}

interface RunningPort {
  protocol: string;
  port: number;
  process: string;
  pid: number;
}

const FirewallRuleRow: React.FC<{ rule: FirewallRule }> = ({ rule }) => {
  return (
    <motion.tr
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
      className="border-b border-sky-100 bg-white hover:bg-sky-50 transition-colors"
    >
      <td className="py-4 px-6 text-sky-800 font-semibold">{rule.id}</td>
      <td className="py-4 px-6 text-sky-700">{rule.protocol}</td>
      <td className="py-4 px-6 text-sky-700">{rule.port}</td>
      <td className="py-4 px-6 text-sky-700">{rule.source}</td>
      <td className="py-4 px-6 text-sky-700">{rule.destination}</td>
      <td className="py-4 px-6">
        <span
          className={`inline-block px-3 py-1 rounded-full text-sm font-medium ${
            rule.action === "ACCEPT"
              ? "bg-green-100 text-green-600"
              : "bg-red-100 text-red-600"
          }`}
        >
          {rule.action}
        </span>
      </td>
    </motion.tr>
  );
};

const RunningPortRow: React.FC<{ port: RunningPort; index: number }> = ({ port, index }) => {
  const { sendAction } = useSocket();

  const handleStopPort = () => {
    sendAction("stop_port", { port: port.port, pid: port.pid });
    toast.info(`Stopping port ${port.port} (PID: ${port.pid})`);
  };

  return (
    <motion.tr
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
      className="border-b border-sky-100 bg-white hover:bg-sky-50 transition-colors"
    >
      <td className="py-4 px-6 text-sky-800 font-semibold">{port.port}</td>
      <td className="py-4 px-6 text-sky-700">{port.protocol}</td>
      <td className="py-4 px-6 text-sky-700">{port.process}</td>
      <td className="py-4 px-6 text-sky-700">{port.pid}</td>
      <td className="py-4 px-6">
        <motion.button
          whileHover={{ scale: 1.1, rotate: 5 }}
          whileTap={{ scale: 0.9 }}
          onClick={handleStopPort}
          className="p-2 bg-red-500 text-white rounded-full hover:bg-red-600 transition-colors"
          title="Stop Port"
        >
          <FaStop size={16} />
        </motion.button>
      </td>
    </motion.tr>
  );
};

const AddPortModal: React.FC<{ onClose: () => void }> = ({ onClose }) => {
  const { sendAction } = useSocket();
  const [port, setPort] = useState("");
  const [protocol, setProtocol] = useState("TCP");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const portNumber = Number(port);
    if (!port || isNaN(portNumber) || portNumber < 1 || portNumber > 65535) {
      toast.error("Please enter a valid port number (1-65535)");
      return;
    }
    sendAction("add_port", { port: portNumber, protocol });
    toast.info(`Adding port ${port} (${protocol})`);
    onClose();
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
      <motion.div
        initial={{ scale: 0.8, y: 50 }}
        animate={{ scale: 1, y: 0 }}
        exit={{ scale: 0.8, y: 50 }}
        className="bg-white rounded-2xl p-6 w-full max-w-md shadow-2xl"
      >
        <div className="flex justify-between items-center mb-4">
          <h4 className="text-xl font-bold text-sky-800">Add New Port</h4>
          <motion.button
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            onClick={onClose}
            className="text-sky-600 hover:text-sky-800"
          >
            <FaTimes size={20} />
          </motion.button>
        </div>
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-sky-700 font-medium mb-1">Port Number</label>
            <input
              type="text"
              value={port}
              onChange={(e) => setPort(e.target.value)}
              className="w-full p-2 border border-sky-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500"
              placeholder="Enter port number (1-65535)"
            />
          </div>
          <div className="mb-4">
            <label className="block text-sky-700 font-medium mb-1">Protocol</label>
            <select
              value={protocol}
              onChange={(e) => setProtocol(e.target.value)}
              className="w-full p-2 border border-sky-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500"
            >
              <option value="TCP">TCP</option>
              <option value="UDP">UDP</option>
            </select>
          </div>
          <div className="flex justify-end space-x-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 bg-gray-300 text-gray-800 rounded-lg hover:bg-gray-400"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-sky-500 text-white rounded-lg hover:bg-sky-600"
            >
              Add Port
            </button>
          </div>
        </form>
      </motion.div>
    </motion.div>
  );
};

const AddFirewallRuleModal: React.FC<{ onClose: () => void }> = ({ onClose }) => {
  const { sendAction } = useSocket();
  const [port, setPort] = useState("");
  const [protocol, setProtocol] = useState("TCP");
  const [source, setSource] = useState("any");
  const [destination, setDestination] = useState("any");
  const [action, setAction] = useState("allow");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const portNumber = Number(port);
    if (!port || isNaN(portNumber) || portNumber < 1 || portNumber > 65535) {
      toast.error("Please enter a valid port number (1-65535)");
      return;
    }
    sendAction("add_rule", {
      port: portNumber,
      protocol,
      source,
      destination,
      action,
    });
    toast.info(`Adding firewall rule for ${port}/${protocol}`);
    onClose();
  };

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
      <motion.div
        initial={{ scale: 0.8, y: 50 }}
        animate={{ scale: 1, y: 0 }}
        exit={{ scale: 0.8, y: 50 }}
        className="bg-white rounded-2xl p-6 w-full max-w-md shadow-2xl"
      >
        <div className="flex justify-between items-center mb-4">
          <h4 className="text-xl font-bold text-sky-800">Add Firewall Rule</h4>
          <motion.button
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            onClick={onClose}
            className="text-sky-600 hover:text-sky-800"
          >
            <FaTimes size={20} />
          </motion.button>
        </div>
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-sky-700 font-medium mb-1">Port Number</label>
            <input
              type="text"
              value={port}
              onChange={(e) => setPort(e.target.value)}
              className="w-full p-2 border border-sky-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500"
              placeholder="Enter port number (1-65535)"
            />
          </div>
          <div className="mb-4">
            <label className="block text-sky-700 font-medium mb-1">Protocol</label>
            <select
              value={protocol}
              onChange={(e) => setProtocol(e.target.value)}
              className="w-full p-2 border border-sky-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500"
            >
              <option value="TCP">TCP</option>
              <option value="UDP">UDP</option>
            </select>
          </div>
          <div className="mb-4">
            <label className="block text-sky-700 font-medium mb-1">Source</label>
            <input
              type="text"
              value={source}
              onChange={(e) => setSource(e.target.value)}
              className="w-full p-2 border border-sky-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500"
              placeholder="e.g., any, 192.168.1.0/24"
            />
          </div>
          <div className="mb-4">
            <label className="block text-sky-700 font-medium mb-1">Destination</label>
            <input
              type="text"
              value={destination}
              onChange={(e) => setDestination(e.target.value)}
              className="w-full p-2 border border-sky-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500"
              placeholder="e.g., any, 192.168.1.100"
            />
          </div>
          <div className="mb-4">
            <label className="block text-sky-700 font-medium mb-1">Action</label>
            <select
              value={action}
              onChange={(e) => setAction(e.target.value)}
              className="w-full p-2 border border-sky-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-sky-500"
            >
              <option value="allow">Allow</option>
              <option value="deny">Deny</option>
            </select>
          </div>
          <div className="flex justify-end space-x-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 bg-gray-300 text-gray-800 rounded-lg hover:bg-gray-400"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-sky-500 text-white rounded-lg hover:bg-sky-600"
            >
              Add Rule
            </button>
          </div>
        </form>
      </motion.div>
    </motion.div>
  );
};

const FirewallAndPortsPage: React.FC = () => {
  const { firewallRules, runningPorts } = useSocket();
  const [showAddPortModal, setShowAddPortModal] = useState(false);
  const [showAddFirewallRuleModal, setShowAddFirewallRuleModal] = useState(false);

  return (
    <div className="container mx-auto px-4">
      <motion.h1
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="text-4xl font-extrabold text-sky-800 mb-8"
      >
        Firewall & Ports
      </motion.h1>
      <div className="mb-12">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-2xl font-bold text-sky-800">Firewall Rules</h2>
          <motion.button
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            onClick={() => setShowAddFirewallRuleModal(true)}
            className="px-4 py-2 bg-sky-500 text-white rounded-lg hover:bg-sky-600 flex items-center"
          >
            <FaPlus className="mr-2" /> Add Firewall Rule
          </motion.button>
        </div>
        {firewallRules ? (
          firewallRules.length > 0 ? (
            <div className="overflow-x-auto bg-white rounded-2xl shadow-xl">
              <table className="w-full table-auto">
                <thead>
                  <tr className="bg-sky-600 text-white">
                    <th className="py-3 px-6 text-left text-sm font-semibold">Rule ID</th>
                    <th className="py-3 px-6 text-left text-sm font-semibold">Protocol</th>
                    <th className="py-3 px-6 text-left text-sm font-semibold">Port</th>
                    <th className="py-3 px-6 text-left text-sm font-semibold">Source</th>
                    <th className="py-3 px-6 text-left text-sm font-semibold">Destination</th>
                    <th className="py-3 px-6 text-left text-sm font-semibold">Action</th>
                  </tr>
                </thead>
                <tbody>
                  {firewallRules.map((rule) => (
                    <FirewallRuleRow key={rule.id} rule={rule} />
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <motion.p
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              className="text-center text-sky-600 text-lg"
            >
              No firewall rules found.
            </motion.p>
          )
        ) : (
          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="text-center text-sky-600 text-lg"
          >
            Loading firewall rules...
          </motion.p>
        )}
      </div>

      <div>
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-2xl font-bold text-sky-800">Running Ports</h2>
          <motion.button
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
            onClick={() => setShowAddPortModal(true)}
            className="px-4 py-2 bg-sky-500 text-white rounded-lg hover:bg-sky-600 flex items-center"
          >
            <FaPlus className="mr-2" /> Add Port
          </motion.button>
        </div>
        {runningPorts ? (
          runningPorts.length > 0 ? (
            <div className="overflow-x-auto bg-white rounded-2xl shadow-xl">
              <table className="w-full table-auto">
                <thead>
                  <tr className="bg-sky-600 text-white">
                    <th className="py-3 px-6 text-left text-sm font-semibold">Port</th>
                    <th className="py-3 px-6 text-left text-sm font-semibold">Protocol</th>
                    <th className="py-3 px-6 text-left text-sm font-semibold">Process</th>
                    <th className="py-3 px-6 text-left text-sm font-semibold">PID</th>
                    <th className="py-3 px-6 text-left text-sm font-semibold">Action</th>
                  </tr>
                </thead>
                <tbody>
                  {runningPorts.map((port, index) => (
                    <RunningPortRow
                      key={`${port.port}-${port.pid}-${port.protocol}-${index}`}
                      port={port}
                      index={index}
                    />
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <motion.p
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              className="text-center text-sky-600 text-lg"
            >
              No running ports found.
            </motion.p>
          )
        ) : (
          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="text-center text-sky-600 text-lg"
          >
            Loading running ports...
          </motion.p>
        )}
      </div>

      <AnimatePresence>
        {showAddPortModal && <AddPortModal onClose={() => setShowAddPortModal(false)} />}
        {showAddFirewallRuleModal && (
          <AddFirewallRuleModal onClose={() => setShowAddFirewallRuleModal(false)} />
        )}
      </AnimatePresence>
    </div>
  );
};

export default FirewallAndPortsPage;