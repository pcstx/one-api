import { Navigate,useLocation } from 'react-router-dom';
import { UserContext } from '../context/User';
import React, { useContext, useState,useEffect } from 'react';

import { history,getCookie,API } from '../helpers';


function PrivateRoute({ children }) {

  const location = useLocation();
  const pathname = location.pathname;
  const [userState, userDispatch] = useContext(UserContext);


  if(getCookie('pushToken') == null ||getCookie('pushToken')==='' ){
    userDispatch({ type: 'logout' });
    localStorage.removeItem('user');
    return <Navigate to='/login' state={{ from: pathname }} />;
  } else if (!localStorage.getItem('user')) {
    //请求接口获取对象
    if(getUserInfo()===0){
      userDispatch({ type: 'logout' });
      localStorage.removeItem('user');
      return <Navigate to='/login' state={{ from: pathname }} />;
    }
  }
  return children;
}

const getUserInfo = async ()=>{
  let res  = await API.get(`/api/user/self`);
  const {success, message, data} = res.data;
  if (success) {
    localStorage.setItem('user',JSON.stringify(data))
    return 1;
  } else {
    return 0;
  }
}

export { PrivateRoute };