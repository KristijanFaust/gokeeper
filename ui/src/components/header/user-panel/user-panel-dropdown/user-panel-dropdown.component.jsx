import './user-panel-dropdown.component.scss'
import {useHistory} from "react-router-dom";
import {useState} from "react";

const signOut = (history, signOutCallback, setSignedOut) => {
  // Since we currently don't store anything authentication related on the backend no API mutations are needed
  localStorage.clear();
  setSignedOut(true);
  signOutCallback(null);
};

const UserPanelDropdown = ({signOutCallback}) => {
  let history = useHistory();
  const [signedOut, setSignedOut] = useState(false)

  if (!signedOut) {
    return (
      <div className='user-panel-dropdown'>
        <span className='title'>USER PANEL</span>
        <span className='action' onClick={() => signOut(history, signOutCallback, setSignedOut)}>SIGN OUT</span>
      </div>
    );
  } else {
    return null
  }
};

export default UserPanelDropdown;