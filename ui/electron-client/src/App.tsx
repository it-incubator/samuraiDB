import React, { ChangeEvent, useEffect, useState } from 'react';
import { v4 as uuidv4 } from 'uuid';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs.tsx';
import SyntaxHighlighter from 'react-syntax-highlighter';

import { github } from 'react-syntax-highlighter/dist/cjs/styles/hljs';


const App = () => {
  const [status, setStatus] = useState<string>('');
  const [payloadData, setPayloadData] = useState<string>(`{
  "id": "123",
  "name": "Miyamoto Musashi",
  "health": 100,
  "attackPower": 50,
  "defensePower": 30,
  "weapon": "Katana"
}`);
  const [getData, setGetData] = useState<string>('53dd6069-3ecc-46c1-ba79-addd1924b997');
  const [response, setResponse] = useState<string>('');
  const [jsonError, setJsonError] = useState<string | null>(null);

  const handleConnect = async () => {
    await window.tcpClient.connectToServer('localhost', 4001);

    window.tcpClient.onData((receivedData: string) => {
      console.log('Received data in React:', receivedData);
      setResponse(receivedData);
    });
  };

  const handleSet = async () => {
    console.log(payloadData);
    await window.tcpClient.sendData({
      type: 'SET',
      payload: JSON.parse(payloadData),
      uuid: uuidv4(),
    });
  };

  const handleGet = async () => {
    await window.tcpClient.sendData({
      type: 'GET',
      payload: { id: getData },
      uuid: uuidv4(),
    });
  };

  const onSetDataInputChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
    setPayloadData(e.target.value);

    try {
      JSON.parse(e.target.value);
      setJsonError(null);
    } catch (error) {
      setJsonError('Invalid JSON');
    }

  };

  useEffect(() => {
    window.tcpClient.onStatus((receivedData: string) => {
      setStatus(receivedData);
    });
  }, []);

  return (
    <div className="flex h-screen bg-gray-100">
      <div className="w-64 bg-white shadow-md">
        <div className="p-4">
          <h2 className="text-xl font-bold mb-4">Database Client</h2>
          <Button onClick={handleConnect} className="w-full mb-2">Connect</Button>
          <Badge variant="outline"
                 className="w-full mb-4">{status}</Badge>
        </div>
      </div>
      <div className="flex-1 p-8 overflow-auto">
        <Tabs defaultValue="set" className="w-full">
          <TabsList className="mb-4">
            <TabsTrigger value="set">Set</TabsTrigger>
            <TabsTrigger value="get">Get</TabsTrigger>
          </TabsList>
          <TabsContent value="set">
            <Card>
              <CardHeader>
                <CardTitle>Set Data</CardTitle>
                <CardDescription>Send data to the server</CardDescription>
              </CardHeader>
              <CardContent>
                <Textarea
                  value={payloadData}
                  onChange={onSetDataInputChange}
                  placeholder="Enter JSON data"
                  className="w-full h-32 mb-4"
                />
                {jsonError && <p className="text-red-500 mb-2">{jsonError}</p>}
                <Button onClick={handleSet} disabled={!!jsonError}>Set Data</Button>
              </CardContent>
            </Card>
          </TabsContent>
          <TabsContent value="get">
            <Card>
              <CardHeader>
                <CardTitle>Get Data</CardTitle>
                <CardDescription>Retrieve data from the server</CardDescription>
              </CardHeader>
              <CardContent>
                <Input
                  value={getData}
                  onChange={(e) => setGetData(e.target.value)}
                  placeholder="Enter ID"
                  className="w-full mb-4"
                />
                <Button onClick={handleGet}>Get Data</Button>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
        <div className="mt-8">
          <h3 className="text-lg font-semibold mb-2">Server Response:</h3>
          <pre className="bg-gray-200 p-4 rounded-md overflow-auto max-h-64">
           {response
             ? <SyntaxHighlighter
               style={github}
               language={'json'}
             >
               {JSON.stringify(JSON.parse(response), null, 2)}
             </SyntaxHighlighter>
             : 'No response received'
           }
          </pre>
        </div>
      </div>
    </div>);
};

export default App;