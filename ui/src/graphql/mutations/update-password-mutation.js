import {gql} from '@apollo/react-hooks';

export default gql`
    mutation UpdatePassword($passwordId: ID!, $name: String!, $password: String!){
        updatePassword(input: {id: $passwordId, name: $name, password: $password}){
            name
            password
        }
    }
`;