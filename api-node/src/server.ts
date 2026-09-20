import "dotenv/config";
import { PrismaPg } from "@prisma/adapter-pg";
import Fastify from "fastify";
import { PrismaClient } from "@prisma/client";

// Inicializa o servidor Fastify e o cliente do Prisma
const fastify = Fastify({ logger: true });
const connectionString = process.env.DATABASE_URL;

if (!connectionString) {
  throw new Error("DATABASE_URL não foi definida.");
}

const adapter = new PrismaPg({ connectionString });
const prisma = new PrismaClient({ adapter });

// Rota 1: Healthcheck (Para saber se a API está viva)
fastify.get("/", async () => {
  return { status: "🛡️ Web3 Shield API Online", version: "1.0.0" };
});

// Rota 2: Buscar as últimas transações analisadas e seus riscos
fastify.get("/risks", async (_request, reply) => {
  try {
    const recentAnalyses = await prisma.risk_analysis.findMany({
      orderBy: {
        analyzed_at: "desc", // Traz as mais recentes primeiro
      },
      take: 10, // Limita a 10 resultados para não sobrecarregar
      include: {
        raw_transactions: true, // Faz um JOIN automático para trazer os dados da transação original!
      },
    });

    const data = recentAnalyses.map((analysis) => ({
      ...analysis,
      analyzed_at: analysis.analyzed_at?.toISOString() ?? null,
      raw_transactions: {
        ...analysis.raw_transactions,
        block_number: analysis.raw_transactions.block_number?.toString() ?? null,
        timestamp: analysis.raw_transactions.timestamp?.toISOString() ?? null,
      },
    }));

    return reply.send({ success: true, data });
  } catch (error) {
    fastify.log.error(error);
    return reply.status(500).send({ success: false, error: "Erro ao buscar dados no banco" });
  }
});

// Rota 3: Consulta de risco para uma carteira específica
fastify.get('/wallet/:address/risk', async (request, reply) => {
  // Extrai o endereço da URL (ex: /wallet/0x123.../risk)
  const { address } = request.params as { address: string };

  try {
    // Busca no banco todas as transações enviadas para este endereço
    const interactions = await prisma.raw_transactions.findMany({
      where: {
        to_address: {
          equals: address,
          mode: 'insensitive', // Evita problemas com hexadecimais em maiúsculo/minúsculo
        },
      },
      include: {
        risk_analysis: true, // Traz o Risk Score calculado pelo Python
      },
    });

    if (interactions.length === 0) {
      return reply.send({
        success: true,
        wallet: address,
        status: 'SAFE_OR_UNKNOWN',
        message: 'Nenhuma interação registrada ou sem risco detectado.',
        interactions: [],
      });
    }

    // Varre as interações para descobrir se a carteira interagiu com contratos de alto risco
    const hasCriticalRisk = interactions.some(
      (tx) => tx.risk_analysis && tx.risk_analysis.risk_score >= 80
    );

    return reply.send({
      success: true,
      wallet: address,
      status: hasCriticalRisk ? 'CRITICAL_RISK' : 'MODERATE_RISK',
      total_interactions: interactions.length,
      interactions: interactions.map((tx) => ({
        hash: tx.tx_hash,
        signature: tx.function_signature,
        risk_score: tx.risk_analysis?.risk_score || 0,
        analyzed_at: tx.risk_analysis?.analyzed_at,
      })),
    });
  } catch (error) {
    fastify.log.error(error);
    return reply.status(500).send({ success: false, error: 'Erro ao analisar a carteira' });
  }
});

// Função de inicialização
const start = async () => {
  try {
    const port = Number(process.env.PORT) || 3000;
    await fastify.listen({ port, host: "0.0.0.0" });
    console.log(`🚀 Servidor rodando em http://localhost:${port}`);
  } catch (err) {
    fastify.log.error(err);
    process.exit(1);
  }
};

start();
