import {useState} from 'react';

import './user-panel-dropdown.component.scss';

const signOut = (signOutCallback, setSignedOut) => {
  // Since we currently don't store anything authentication related on the backend no API mutations are needed
  localStorage.clear();
  setSignedOut(true);
  signOutCallback('');
};

const UserPanelDropdown = ({signOutCallback}) => {
  const [signedOut, setSignedOut] = useState(false);

  if (!signedOut) {
    return (
      <div className='user-panel-dropdown'>
        <span className='title'>USER PANEL</span>
        <span className='action' onClick={() => signOut(signOutCallback, setSignedOut)}>SIGN OUT</span>
      </div>
    );
  } else {
    return null;
  }
};

export default UserPanelDropdown;