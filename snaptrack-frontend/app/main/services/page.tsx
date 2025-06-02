"use client";

import React, { useState } from "react";
import { toast, ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { motion, AnimatePresence } from "framer-motion";
import { FaPlay, FaStop, FaSync, FaFileAlt, FaTimes } from "react-icons/fa";
import { useSocket } from "@/app/context/SocketContext";

interface ServiceInfo {
  name: string;
  status: string;
  uptime: string;
  memory: string;
  version: string;
}

const ServiceRow: React.FC<{ service: ServiceInfo }> = ({ service }) => {
  const { sendAction, logs } = useSocket();
  const [showLogs, setShowLogs] = useState<boolean>(false);

  const handleAction = (action: string) => {
    sendAction(action, service.name);
    toast.success(`${action.charAt(0).toUpperCase() + action.slice(1)} initiated for ${service.name}`, {
      position: "top-right",
      autoClose: 2000,
    });
  };

  return (
    <>
      <motion.tr
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
        className="border-b border-sky-100 bg-white hover:bg-sky-50 transition-colors"
      >
        <td className="py-4 px-6 text-sky-800 font-semibold">{service.name}</td>
        <td className="py-4 px-6">
          <span
            className={`inline-block px-3 py-1 rounded-full text-sm font-medium ${
              service.status === "active"
                ? "bg-green-100 text-green-600"
                : "bg-red-100 text-red-600"
            }`}
          >
            {service.status}
          </span>
        </td>
        <td className="py-4 px-6 text-sky-700">{service.version}</td>
        <td className="py-4 px-6 text-sky-700">{service.uptime}</td>
        <td className="py-4 px-6 text-sky-700">{service.memory}</td>
        <td className="py-4 px-6 flex space-x-3">
          <motion.button
            whileHover={{ scale: 1.1, rotate: 5 }}
            whileTap={{ scale: 0.9 }}
            onClick={() => handleAction("start")}
            className="p-2 bg-green-500 text-white rounded-full hover:bg-green-600 transition-colors"
            title="Start Service"
          >
            <FaPlay size={16} />
          </motion.button>
          <motion.button
            whileHover={{ scale: 1.1, rotate: 5 }}
            whileTap={{ scale: 0.9 }}
            onClick={() => handleAction("stop")}
            className="p-2 bg-red-500 text-white rounded-full hover:bg-red-600 transition-colors"
            title="Stop Service"
          >
            <FaStop size={16} />
          </motion.button>
          <motion.button
            whileHover={{ scale: 1.1, rotate: 5 }}
            whileTap={{ scale: 0.9 }}
            onClick={() => handleAction("restart")}
            className="p-2 bg-blue-500 text-white rounded-full hover:bg-blue-600 transition-colors"
            title="Restart Service"
          >
            <FaSync size={16} />
          </motion.button>
          <motion.button
            whileHover={{ scale: 1.1, rotate: 5 }}
            whileTap={{ scale: 0.9 }}
            onClick={() => {
              handleAction("logs");
              setShowLogs(true);
            }}
            className="p-2 bg-sky-500 text-white rounded-full hover:bg-sky-600 transition-colors"
            title="View Logs"
          >
            <FaFileAlt size={16} />
          </motion.button>
        </td>
      </motion.tr>

      <AnimatePresence>
        {showLogs && logs[service.name] && (
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
              className="bg-white rounded-2xl p-6 w-full max-w-2xl max-h-[80vh] flex flex-col shadow-2xl"
            >
              <div className="flex justify-between items-center mb-4">
                <h4 className="text-xl font-bold text-sky-800">Logs for {service.name}</h4>
                <motion.button
                  whileHover={{ scale: 1.1 }}
                  whileTap={{ scale: 0.9 }}
                  onClick={() => setShowLogs(false)}
                  className="text-sky-600 hover:text-sky-800"
                >
                  <FaTimes size={20} />
                </motion.button>
              </div>
              <div className="bg-sky-50 p-4 rounded-lg max-h-96 overflow-y-auto text-sm text-sky-700">
                {logs[service.name].map((log: string, index: number) => (
                  <motion.p
                    key={index}
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: index * 0.05 }}
                    className="py-1 border-b border-sky-100 last:border-b-0"
                  >
                    {log}
                  </motion.p>
                ))}
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </>
  );
};

const ServicesPage: React.FC = () => {
  const { services } = useSocket();

  return (
    <div >
      <motion.h1
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="text-4xl font-extrabold text-sky-800 mb-8 text-start"
      >
        System Services Dashboard
      </motion.h1>
      <div className="container mx-auto">
        {services ? (
          services.length > 0 ? (
            <div className="overflow-x-auto bg-white rounded-2xl shadow-xl">
              <table className="w-full table-auto">
                <thead>
                  <tr className="bg-sky-600 text-white">
                    <th className="py-3 px-6 text-left text-sm font-semibold">Service</th>
                    <th className="py-3 px-6 text-left text-sm font-semibold">Status</th>
                    <th className="py-3 px-6 text-left text-sm font-semibold">Version</th>
                    <th className="py-3 px-6 text-left text-sm font-semibold">Uptime</th>
                    <th className="py-3 px-6 text-left text-sm font-semibold">Memory</th>
                    <th className="py-3 px-6 text-left text-sm font-semibold">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {services.map((service: ServiceInfo) => (
                    <ServiceRow key={service.name} service={service} />
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
              No services found.
            </motion.p>
          )
        ) : (
          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="text-center text-sky-600 text-lg"
          >
            Loading services...
          </motion.p>
        )}
      </div>
      <ToastContainer position="top-right" autoClose={2000} theme="colored" />
    </div>
  );
};

export default ServicesPage;