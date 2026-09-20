import "dotenv/config";
import { PrismaPg } from "@prisma/adapter-pg";
import Fastify from "fastify";
import { PrismaClient } from "@prisma/client";

// Serviços HTTP e acesso ao PostgreSQL.
const fastify = Fastify({ logger: true });
const connectionString = process.env.DATABASE_URL;

if (!connectionString) {
  throw new Error("DATABASE_URL não foi definida.");
}

const adapter = new PrismaPg({ connectionString });
const prisma = new PrismaClient({ adapter });

// Endpoint de disponibilidade da API.
fastify.get("/", async () => {
  return { status: "🛡️ Web3 Shield API Online", version: "1.0.0" };
});

// Retorna as análises de risco mais recentes com a transação relacionada.
fastify.get("/risks", async (_request, reply) => {
  try {
    const recentAnalyses = await prisma.risk_analysis.findMany({
      orderBy: {
        analyzed_at: "desc",
      },
      take: 10,
      include: {
        raw_transactions: true,
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

// Consolida o risco das interações associadas a uma carteira.
fastify.get('/wallet/:address/risk', async (request, reply) => {
  const { address } = request.params as { address: string };

  try {
    const interactions = await prisma.raw_transactions.findMany({
      where: {
        to_address: {
          equals: address,
          mode: 'insensitive',
        },
      },
      include: {
        risk_analysis: true,
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

// Inicia o servidor na porta configurada pelo ambiente.
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
