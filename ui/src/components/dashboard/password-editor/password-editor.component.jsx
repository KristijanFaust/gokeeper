import PasswordEdit from './password-edit/password-edit.component';

import './password-editor.component.scss';

const PasswordEditor = ({passwords, authenticationExpiredCallback}) => {
  return (
    <div className='password-editor'>
      {passwords.map((password) => (
        <PasswordEdit key={password.id} passwordEntry={password} authenticationExpiredCallback={authenticationExpiredCallback}/>
      ))}
    </div>
  );
};

export default PasswordEditor;