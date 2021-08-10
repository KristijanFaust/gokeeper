import React from "react";

import Logo from '../logo/logo.component';
import HeaderLinks from "./header-links/header-links.component";
import UserPanel from "./user-panel/user-panel.component";

import './header.styles.scss';

const Header = ({navigationLinks}) => (
    <div className='header'>
      <Logo />
      <HeaderLinks navigationLinks={navigationLinks} />
      <UserPanel />
    </div>
);

export default Header;