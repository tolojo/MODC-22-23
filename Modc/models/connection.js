const pg = require('pg');

let pool = new pg.Pool({
    user: 'postgres',/*  mudar para o nome */
    host: 'localhost',
    password: '1234',/* mudar a pass */
    port: 5432,
    ssl: false
});

pool.connect((err, client, done) => {
    if (err) throw err;
  
    client.query('CREATE DATABASE modc', (error, result) => {
      done();
      if (error) {
        console.error('Erro ao criar a base de dados:', error);
      } else {
        console.log('Base de dados criada!'); 
        pool = new pg.Pool({
            user: 'postgres',/*  mudar para o nome */
            host: 'localhost',
            database: 'modc',
            password: '1234',/* mudar a pass */
            port: 5432,
            ssl: false
        });
        criarTabela();
      }
    });
  });
  
  function criarTabela() {
    pool.query(`
      CREATE TABLE IF NOT EXISTS utilizadores (
        use_id serial,
        use_nome varchar(255),
        use_email varchar(255),
        use_password varchar(255),
        primary Key (use_id));`, (error, result) => {
      if (error) {
        console.error('Erro ao criar a tabela:', error);
      } else {
        console.log('Tabela criada com sucesso!');
      }
  
      pool.end();
    });
  }

  const connectionString = "postgres://postgres:1234@localhost:5432/modc"/*  mudar tambem o nome e a pass*/
    const Pool = pg.Pool
    pool = new Pool({
    connectionString,
    max: 10,
    
})

  
module.exports = pool;