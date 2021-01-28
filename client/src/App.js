import './App.css';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';
import UserProfile from './components/pages/Login'
import Register from './components/pages/Register'
import UserDashboard from './components/pages/UserDashboard'

function App() {
  return (
    <>
    <Router>
      <Switch>
          <Route exact path='/' component={UserProfile}/>
          <Route exact path='/register' component={Register}/>
          <Route exact path='/dashboard' component={UserDashboard}/>
      </Switch>
    </Router>
    </>

  );
}

export default App;
