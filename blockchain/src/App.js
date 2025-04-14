// React frontend code

import React, { useEffect, useState } from 'react';
import axios from 'axios';

export default function BlockchainApp() {
  const [wallet, setWallet] = useState('');
  const [receiver, setReceiver] = useState('');
  const [amount, setAmount] = useState('');
  const [balance, setBalance] = useState(0);
  const [transactions, setTransactions] = useState([]);
  const [blocks, setBlocks] = useState([]);
  const [message, setMessage] = useState('');

  useEffect(() => {
    if (wallet) {
      fetchBalance();
    }
  }, [wallet]);

  const createWallet = async () => {
    const res = await axios.get('http://localhost:8080/wallet/new');
    setWallet(res.data.public);
    setMessage('Wallet created!');
  };

  const fetchBalance = async () => {
    const res = await axios.get(`http://localhost:8080/balance/${wallet}`);
    setBalance(res.data.balance);
    setTransactions(res.data.transactions);
  };

  const sendTransaction = async () => {
    try {
      await axios.post('http://localhost:8080/transaction', {
        sender: wallet,
        receiver,
        amount: parseFloat(amount),
        signature: 'placeholder-signature'
      });
      setMessage('Transaction sent! Pending mining...');
      setReceiver('');
      setAmount('');
    } catch (err) {
      setMessage('Error: Transaction failed');
    }
  };

  const mineBlock = async () => {
    const res = await axios.get('http://localhost:8080/mine');
    setMessage(`Block mined! Hash: ${res.data.hash}`);
    fetchBlocks();
    fetchBalance();
  };

  const fetchBlocks = async () => {
    const res = await axios.get('http://localhost:8080/blocks');
    setBlocks(res.data);
  };

  useEffect(() => {
    fetchBlocks();
  }, []);

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-4">Blockchain Payment System</h1>

      {!wallet && <button className="bg-blue-500 text-white px-4 py-2 rounded" onClick={createWallet}>Create Wallet</button>}

      {wallet && (
        <div className="space-y-4">
          <p className="text-sm">Your Wallet Address:</p>
          <code className="block p-2 bg-gray-100 rounded break-words">{wallet}</code>

          <p>Balance: <strong>{balance}</strong> Tokens</p>

          <div className="space-y-2">
            <input
              className="border p-2 w-full"
              placeholder="Receiver address"
              value={receiver}
              onChange={e => setReceiver(e.target.value)}
            />
            <input
              className="border p-2 w-full"
              placeholder="Amount"
              type="number"
              value={amount}
              onChange={e => setAmount(e.target.value)}
            />
            <button className="bg-green-500 text-white px-4 py-2 rounded" onClick={sendTransaction}>
              Send Transaction
            </button>
          </div>

          <button className="bg-yellow-500 text-white px-4 py-2 mt-4 rounded" onClick={mineBlock}>Mine Block</button>

          {message && <p className="mt-4 text-blue-600 font-semibold">{message}</p>}

          <div className="mt-6">
            <h2 className="text-xl font-bold mb-2">Your Transactions</h2>
            <ul className="space-y-2">
              {transactions?.map((tx, idx) => (
                <li key={idx} className="bg-white p-2 border rounded">
                  From: {tx.sender.slice(0, 10)}... → {tx.receiver.slice(0, 10)}... | {tx.amount} tokens
                </li>
              ))}
            </ul>
          </div>

          <div className="mt-6">
            <h2 className="text-xl font-bold mb-2">Blockchain Explorer</h2>
            <ul className="space-y-4">
              {blocks.map((block, i) => (
                <li key={i} className="bg-gray-50 p-4 rounded shadow">
                  <p className="font-semibold">Block #{block.index}</p>
                  <p>Hash: <code className="break-words">{block.hash}</code></p>
                  {/* <p>Transactions: {block.transactions.length}</p> */}
                </li>
              ))}
            </ul>
          </div>

        </div>
      )}
    </div>
  );
}
