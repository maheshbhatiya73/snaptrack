"use client";
import React from 'react';
import FirewallList from './FirewallList';
import PortsList from './PortsList';

const App = () => {
  return (
    <div className="">
      <FirewallList />
      <PortsList />
    </div>
  );
};

export default App;