import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import './env.js';
import App from './components/App';
import reportWebVitals from './reportWebVitals';
import { GoogleOAuthProvider } from "@react-oauth/google"
import { BrowserRouter } from "react-router";

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);
root.render(
  <GoogleOAuthProvider clientId='791175547804-b8ti91j4mccc6gin1193oaslpm38adla.apps.googleusercontent.com'>
    <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
    </React.StrictMode>
  </GoogleOAuthProvider>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
