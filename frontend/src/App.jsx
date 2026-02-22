import React from 'react';
import ClickButton from './components/ClickButton';
import 'antd/dist/reset.css';

const App = () => {
  const handleClick = () => {
    alert('Button clicked!');
  };
  return (
<div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
  <ClickButton />
</div>
  );
};

export default App;
