import React from 'react';

export default function InvalidCredentials ({ message, color }) {
  const styles = {
    paddingTop: '5px',
    color: color,
    textAlign: 'center',
  };
  return (
    <div style={styles}>
      {message}
    </div>
  );
};