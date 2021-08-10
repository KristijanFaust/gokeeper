import './button.styles.scss';

const Button = ({children, disabled}) => (
  <button className={disabled ? 'button disabled' : 'button'} disabled={disabled}>
    {children}
  </button>
);

export default Button;
