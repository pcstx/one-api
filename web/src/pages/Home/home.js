import React, { useContext, useEffect, useState } from 'react';
import { Card, Grid, Header, Segment, Container,Button } from 'semantic-ui-react';
import { API, showError, showNotice } from '../../helpers';
import { StatusContext } from '../../context/Status';
import { marked } from 'marked';

const NewHome = () => {
  const [statusState, statusDispatch] = useContext(StatusContext);
  const [homePageContentLoaded, setHomePageContentLoaded] = useState(false);
  const [homePageContent, setHomePageContent] = useState('');

  const displayNotice = async () => {
    const res = await API.get('/api/notice');
    const { success, message, data } = res.data;
    if (success) {
      let oldNotice = localStorage.getItem('notice');
      if (data !== oldNotice && data !== '') {
        showNotice(data);
        localStorage.setItem('notice', data);
      }
    } else {
      showError(message);
    }
  };
 

  useEffect(() => {
    displayNotice().then(); 
  }, []);
  return (
    <>
      {
            <Container textAlign="center" style={{margin:"100px 0 100px 0",width:"755px"}}>
                <p style={{textAlign:'left',lineHeight:'30px'}}>
                    <strong style={{fontSize:'46px'}}>破壳AI</strong>
                </p>
                <p style={{textAlign:'left',lineHeight:'30px'}}>
                    <strong style={{fontSize:'46px'}}>生成式AI应用开发工具平台</strong>
                </p>
                <Container  style={{textAlign:'left',fontSize:'18px',lineHeight:'30px'}}>
                整合AI的能力，将其应用范围拓展到更多领域，使AI技术普惠于广大人群。教育、政策制定和提高可访问性等措施将推动AI的全面发展，为社会带来积极的变革和发展机遇。
                </Container>
                <div>
                    <Button class="ui button" size="big" color="blue" style={{margin:'60px',backgroundColor:'#5680ff'}}>开始使用</Button>
                </div>
            </Container>
      }

    </>
  );
};

export default NewHome;
