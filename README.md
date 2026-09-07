# NAO SEI SO TAVA AFIM

Não tem um objetivo específico; estou apenas adicionando coisas que acho legais neste repositório e usando-o para estudo e pesquisa.

# DATA: 07/09/2026:
  ## OBJETIVO: 
  Implementar um pequeno mock de criação de usuario. Sem banco, na memoria, vou experimentar o conceito de portas e adaptadores(arquitetura hexagonal).

  O conceito é um tanto complexo no sentido que ele depende de boas abstraçoes de interface pra funcionar corretamente. 
  Basicamente, ele definir uma esquema parecido com uma cebola (camadas), onde nossas regras de negocio e estruturas mais "puras" 
  estao localizadas bem no centro dessa cebola, e nas camadas mais externas temos o que seria a "fronteira" com o mundo externo. 

  ![arquitetura haxagonal](/assets/arquitetura-hexagonal.svg) 

  O conceito central aqui é que nossas classes de dominio nao
   nao tem nenhum conhecimento sobre a origem ou o destino dos dados que recebe,
   focando na mais pura logica, elas nao dependem de infraestrura, tecnologia ou sistema externo, 
   apenas se comunicando com o externo atraves das portas de entrada e saida.

   ![adaptadores e portas](/assets/hex-ports-adapters.svg)
   ## Portas
   Portas sao a forma como as classes de dominio se comunicam com o mundo externo.
   O termo "porta" se refere as interfaces usadas pra se comunicar com as classes de dominio.
   
   existem dois tipos de portas:
   - **Portas de Entrada**: Sao usadas pra se comunicar de fora pra dentro, é como o sistema declara os serviçoes forncidos ao mundo exterior.
   - **Portas de Saida**: Usados pra se comunicar de dentro pra fora, quando uma classe de dominio precisa de uma servico externo. Ela declara as dependencias 
   - do mundo exterior que o sistema necessita pra funcionar.
   
   ### [Adaptadores](https://refactoring.guru/pt-br/design-patterns/adapter)
   ![adapter-png](/assets/adapter-pt-br.png)
   
   Adaptadores sao um padrao de projeto em que, atraves de uma interface,
   eu permito que objetos incompativeis consigam trabalhar em conjunto.
   Aqui, o padrao Adapter permite que o sistema consiga interagir com
   o mundo externo sem a necessidade de saber quem esta "por tras"
   da interface de antemao.
   O consumidor que precisar interagir com o mundo externo nao precisa
   conhecer os detalhes da implementacao concreta, apenas a porta
   definida pela aplicacao. Isso facilita trocar tecnologias,
   testar comportamentos e manter a logica de dominio isolada.
