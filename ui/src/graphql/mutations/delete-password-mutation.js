import {gql} from '@apollo/react-hooks';

export default gql`
    mutation DeletePassword($passwordId: ID!){
        deletePassword(input: $passwordId)
    }
`;