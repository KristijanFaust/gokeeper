import Logo from './logo/logo.component';
import HeaderLinks from './header-links/header-links.component';
import UserPanel from './user-panel/user-panel.component';

import './header.styles.scss';

const Header = ({navigationLinks, signOutCallback}) => (
    <div className='header'>
      <Logo />
      <HeaderLinks navigationLinks={navigationLinks} />
      <UserPanel signOutCallback={signOutCallback} />
    </div>
);

export default Header;