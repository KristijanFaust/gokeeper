import {Link} from 'react-router-dom';
import {useEffect, useState} from 'react';

import './header-links.styles.scss';

const HeaderLinks = ({navigationLinks}) => {
  const [authenticationToken, setAuthenticationToken] = useState(localStorage.getItem('authenticationToken'));

  useEffect(() => {
    setAuthenticationToken(localStorage.getItem('authenticationToken'));
  }, [authenticationToken]);

  return (
    <div className='header-link-group'>
      {Object.keys(navigationLinks).map((link, index) => (
        <Link key={index} className='header-link' to={'/' + navigationLinks[link]}>
          <span>{link}</span>
        </Link>
        )
      )}
    </div>
  );
};

export default HeaderLinks;
