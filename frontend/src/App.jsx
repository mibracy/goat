import './App.css'
import Layout from './components/Layout'
import Login from './components/Login'
import Admin from './components/Admin'
import Agent from './components/Agent'
import Customer from './components/Customer'

function App() {
  return (
    <Layout>
      {({ activeModel, setApiResponse }) => (
        <>
          {activeModel === 'login' && <Login setApiResponse={setApiResponse} />}
          {activeModel === 'admin' && <Admin activeModel={activeModel} setApiResponse={setApiResponse} />}
          {activeModel === 'agent' && <Agent activeModel={activeModel} setApiResponse={setApiResponse} />}
          {activeModel === 'customer' && <Customer activeModel={activeModel} setApiResponse={setApiResponse} />}
        </>
      )}
    </Layout>
  )
}

export default App
