import { Navigate } from 'react-router-dom';

import { history,getCookie,API } from '../helpers';


function PrivateRoute({ children }) {
  if(getCookie('pushToken') == null){
    return <Navigate to='/login' state={{ from: history.location }} />;
  } else if (!localStorage.getItem('user')) {
    //请求接口获取对象
    if(getUserInfo()===0){
      return <Navigate to='/login' state={{ from: history.location }} />;
    }
  }
  return children;
}

const getUserInfo = async ()=>{
  let res  = await API.get(`/api/user/self`);
  const {success, message, data} = res.data;
  if (success) {
    localStorage.setItem('user',data)
    return 1;
  } else {
    return 0;
  }
}

export { PrivateRoute };