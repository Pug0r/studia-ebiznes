import React, { useState } from 'react'
import axios from 'axios'

const API = 'http://localhost:8080'

export default function App() {
  const oauthSuccess = new URLSearchParams(window.location.search).get('oauth') === 'ok'
  const [mode, setMode] = useState('login')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [message, setMessage] = useState('')

  async function submit(e) {
    e.preventDefault()
    setMessage('')
    try {
      const url = mode === 'login' ? '/login' : '/register'
      const res = await axios.post(API + url, { email, password })
      setMessage(res.data?.status === 'ok' ? 'ok' : 'server response')
    } catch (err) {
      setMessage(err.response?.data || 'error')
    }
  }

  return (
    <main>
      <h1>Lab08 OAuth</h1>
      <a href="http://localhost:8080/auth/google">
        Login with Google
      </a>
      <button onClick={() => setMode('login')} disabled={mode === 'login'}>
        Login
      </button>
      <button onClick={() => setMode('register')} disabled={mode === 'register'}>
        Register
      </button>
      <form onSubmit={submit}>
        <label>
          Email
          <input value={email} onChange={(e) => setEmail(e.target.value)} />
        </label>
        <label>
          Password
          <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        </label>
        <button type="submit">Submit</button>
      </form>
      {oauthSuccess && <p style={{ color: 'green' }}>OAuth login successful</p>}
      {message && <p>{String(message)}</p>}
    </main>
  )
}
