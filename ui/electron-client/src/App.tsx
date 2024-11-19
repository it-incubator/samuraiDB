import {  useState } from 'react';
import './App.css'

import {v4 } from 'uuid'

const App = () => {
  const [status, setStatus] = useState<string>('');
  const [payloadData, setPayloadPayloadData] = useState<string>(`{
  "id": "123",
  "name": "Miyamoto Musashi",
  "health": 100,
  "attackPower": 50,
  "defensePower": 30,
  "weapon": "Katana"
}`);
  const [getData, setGetData] = useState<string>('53dd6069-3ecc-46c1-ba79-addd1924b997');
  const [response, setResponse] = useState<string>('');


  const handleConnect = async  () => {
    await window.tcpClient.connectToServer('localhost', 4001);

    window.tcpClient.onData((receivedData: string) => {
      console.log('Received data in React:', receivedData);
      setResponse(receivedData);
    });

    setStatus('Connected to server');
  };

  const handleSet = async () => {
    await window.tcpClient.sendData({ type: 'SET', payload: JSON.parse(payloadData), uuid: v4() })
  };

  const handleGet = async () => {
    await window.tcpClient.sendData({ type: 'GET', payload: { id: getData }, uuid: v4() })
  };


  return (
    <div>
      <h1>SamuraiDB Interface</h1>
      <div>
        <button onClick={handleConnect}>Connect</button>{" "}
        <span>Status: {status}</span>
      </div>

      <div>
        <input
          type="text"
          value={payloadData}
          onChange={(e) => setPayloadPayloadData(e.target.value)}
          placeholder="Data to set"
        />
        <button onClick={handleSet}>Set</button>
      </div>
      <div>
        <input
          type="text"
          value={getData}
          onChange={(e) => setGetData(e.target.value)}
          placeholder="Data to set"
        />
        <button onClick={handleGet}>Get</button>
      </div>

      <div>Response:
      <div>
        <textarea value={response} rows={20} cols={60}/>
      </div>
      </div>
    </div>
  );
};

export default App
