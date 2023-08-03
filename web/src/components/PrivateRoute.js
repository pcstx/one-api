import { Navigate } from 'react-router-dom';

import { history,getCookie } from '../helpers';


function PrivateRoute({ children }) {
  if(getCookie('pushToken') == null){
    return <Navigate to='/login' state={{ from: history.location }} />;
  }
  return children;
}

export { PrivateRoute };