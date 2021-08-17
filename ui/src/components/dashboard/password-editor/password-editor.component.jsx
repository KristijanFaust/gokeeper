import PasswordEdit from './password-edit/password-edit.component';

import './password-editor.component.scss';

const PasswordEditor = ({passwords}) => {
  return (
    <div className='password-editor'>
      {passwords.map((password) => (
        <PasswordEdit key={password.id} passwordEntry={password} />
      ))}
    </div>
  );
};

export default PasswordEditor;