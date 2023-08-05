import React, { useContext, useState,useEffect } from 'react';
import { Link, useNavigate,useLocation } from 'react-router-dom';
import { UserContext } from '../context/User';

import { Button, Container, Dropdown, Grid, Icon, Menu, Segment } from 'semantic-ui-react';
import { API, getLogo, getSystemName, isAdmin, isMobile, showSuccess } from '../helpers';
import '../index.css';

// Header Buttons
let headerButtons = [
  {
    name: '首页',
    to: '/',
    icon: 'home'
  },
  {
    name: '聊天',
    to: '/chat',
    icon: 'comments'
  },
  {
    name: '渠道',
    to: '/channel',
    icon: 'sitemap',
    admin: true
  },
  {
    name: '令牌',
    to: '/token',
    icon: 'key'
  },
  {
    name: '兑换',
    to: '/redemption',
    icon: 'dollar sign',
    admin: true
  },
  {
    name: '充值',
    to: '/recharge',
    icon: 'cart'
  },
  {
    name: '用户',
    to: '/user',
    icon: 'user',
    admin: true
  },
  {
    name: '日志',
    to: '/log',
    icon: 'book'
  },
  {
    name: '设置',
    to: '/setting',
    icon: 'setting',
    admin: true
  },
  {
    name: '使用说明',
    to: '/about',
    icon: 'info circle'
  }
];

// if (localStorage.getItem('chat_link')) {
//   headerButtons.splice(1, 0, {
//     name: '聊天',
//     to: '/chat',
//     icon: 'comments'
//   });
// }
 
const Header = () => {
  const [userState, userDispatch] = useContext(UserContext);
  let navigate = useNavigate();

  const [showSidebar, setShowSidebar] = useState(false);
  const [activeItem,setActiveItem] = useState("首页");
  const systemName = getSystemName();
  const logo = getLogo();
  const location = useLocation();

  useEffect(() => {
    // 在每次 URL 地址变化时更新 activeItem 状态
    const pathname = location.pathname;
    const activeButton = headerButtons.find((button) => pathname===button.to);

    if (activeButton) {
      setActiveItem(activeButton.name);
    }
  }, [location.pathname]);

  async function logout() {
    setShowSidebar(false);
    await API.get('/api/user/wechatLogout');
    showSuccess('注销成功!');
    userDispatch({ type: 'logout' });
    localStorage.removeItem('user');
    navigate('/');
  }

  const toggleSidebar = () => {
    setShowSidebar(!showSidebar);
  };

  function handleItemClick(name) {
    setActiveItem(name)
  }

  const renderButtons = (isMobile) => {
    return headerButtons.map((button) => {
      if (button.admin && !isAdmin()) return <></>;
      if (isMobile) {
        return (
          <Menu.Item
            onClick={() => {
              navigate(button.to);
              setShowSidebar(false);
            }}
          >
            {button.name}
          </Menu.Item>
        );
      }
      return (
        <Menu.Item key={button.name} as={Link} to={button.to} active={activeItem === button.name} onClick={() => handleItemClick(button.name)}>
          {/* <Icon name={button.icon} /> */}
          {button.name}
        </Menu.Item>
      );
    });
  };

  if (isMobile()) {
    return (
      <>
        <Menu
          borderless
          size='large'
          style={
            showSidebar
              ? {
                borderBottom: 'none',
                marginBottom: '0',
                borderTop: 'none',
                height: '51px'
              }
              : { border: 'none', height: '52px' }
          }
        >
          <Container>
            <Menu.Item as={Link} to='/'>
              <img
                src={logo}
                alt='logo'
                style={{ marginRight: '0.75em' }}
              />
              <div style={{ fontSize: '20px' }}>
                <b>{systemName}</b>
                <p style={{fontSize: '15px'}}>开放平台</p>
              </div>
            </Menu.Item>
            <Menu.Menu position='right'>
              <Menu.Item onClick={toggleSidebar}>
                <Icon name={showSidebar ? 'close' : 'sidebar'} />
              </Menu.Item>
            </Menu.Menu>
          </Container>
        </Menu>
        {showSidebar ? (
          <Segment style={{ marginTop: 0, borderTop: '0' }}>
            <Menu secondary vertical style={{ width: '100%', margin: 0 }}>
              {renderButtons(true)}
              <Menu.Item>
                {userState.user ? (
                  <Button onClick={logout}>注销</Button>
                ) : (
                  <>
                    <Button
                      onClick={() => {
                        setShowSidebar(false);
                        navigate('/login');
                      }}
                    >
                      登录
                    </Button>
                    {/* <Button
                      onClick={() => {
                        setShowSidebar(false);
                        navigate('/register');
                      }}
                    >
                      注册
                    </Button> */}
                  </>
                )}
              </Menu.Item>
            </Menu>
          </Segment>
        ) : (
          <></>
        )}
      </>
    );
  }

  return (
    <>
      <Menu borderless style={{ border: 'none' }}>
        <Container>
          <Menu.Item as={Link} to='/' className={'hide-on-mobile'}>
            <img src={logo} alt='logo' style={{ marginRight: '0.75em' }} />
            <div style={{ fontSize: '20px' }}>
              <b>{systemName}</b>
              <p style={{fontSize: '15px'}}>开放平台</p>
            </div>
          </Menu.Item>
          <Container className={'div-flex'}> 
          {renderButtons(false)}
          </Container>
          <Menu.Menu position='right'>
            {userState.user ? (
              <Dropdown
                text={userState.user.display_name}
                pointing
                className='link item'
              >
                <Dropdown.Menu>
                  <Dropdown.Item onClick={logout}>注销</Dropdown.Item>
                </Dropdown.Menu>
              </Dropdown>
            ) : (
              <Menu.Item
                name='登录'
                as={Link}
                to='/login'
                className='btn btn-link'
              />
            )}
          </Menu.Menu>
        </Container>
      </Menu>
    </>
  );
};

export default Header;
